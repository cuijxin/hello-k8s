package mongo

import (
	"context"
	"fmt"
	"hello-k8s/pkg/model/common"
	"hello-k8s/pkg/storage/database"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *MongoDB) Init(option database.DBInitOptions) (err error) {
	m.client, err = mongo.NewClient(
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/",
			option.User, option.Password, option.Address)),
	)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err = m.client.Connect(ctx); err != nil {
		return
	}

	m.collection = m.client.Database(viper.GetString("constants.addons_database")).Collection(viper.GetString("constants.addons_collection"))
	return
}

func (m *MongoDB) Store(app *common.AtomApplication) (err error) {
	_, err = m.collection.InsertOne(context.Background(), app)
	return
}

func (m *MongoDB) Get(option database.RecordOptions) (*common.AtomApplication, error) {
	var result common.AtomApplication
	filter := bson.M{"name": option.Name, "namespace": option.Namespace, "clusterid": option.ClusterID, "type": option.Type}
	err := m.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MongoDB) ListAddons(option database.RecordOptions) ([]common.AtomApplication, error) {
	filter := bson.M{"namespace": option.Namespace, "clusterid": option.ClusterID, "type": bson.M{"$ne": option.Type}}
	result, err := getList(m, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDB) List(option database.RecordOptions) ([]common.AtomApplication, error) {
	filter := bson.M{"namespace": option.Namespace, "clusterid": option.ClusterID, "type": option.Type}
	result, err := getList(m, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getList(m *MongoDB, filter bson.M) ([]common.AtomApplication, error) {
	cur, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	appsInfo := make([]common.AtomApplication, 0)
	for cur.Next(context.Background()) {
		var result common.AtomApplication
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		appsInfo = append(appsInfo, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return appsInfo, nil
}

func (m *MongoDB) Exist(option database.RecordOptions) (bool, error) {
	var result bson.M
	err := m.collection.FindOne(context.Background(), bson.M{"name": option.Name, "namespace": option.Namespace, "clusterid": option.ClusterID, "type": option.Type}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func (m *MongoDB) Delete(option database.RecordOptions) error {
	filter := bson.M{"name": option.Name, "namespace": option.Namespace, "clusterid": option.ClusterID, "type": option.Type}
	_, err := m.collection.DeleteOne(context.Background(), filter)
	return err
}

func (m *MongoDB) Update(app *common.AtomApplication) error {
	selector := bson.M{"name": app.Name, "namespace": app.Namespace, "clusterid": app.ClusterID, "type": app.Type}
	updateData := bson.M{"$set": app}
	_, err := m.collection.UpdateOne(context.Background(), selector, updateData)
	return err
}
