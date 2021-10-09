package main

import (
    "fmt"
    "net/url"
    "time"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/gorilla/mux"	
)

type Users struct {
    ID        primitive.ObjectID  `json:"ID,omitempty" bson:"ID,omitempty"`
    Name      string `json:"Name,omitempty" bson:"ID,omitempty"`
    Email     string `json:"Email,omitempty" bson:"ID,omitempty"`
    Password  string `json:"Password,omitempty" bson:"ID,omitempty"`
}

type Posts struct {
    ID        primitive.ObjectID  `json:"ID,omitempty" bson:"ID,omitempty"`
    Caption      string `json:"Caption,omitempty" bson:"ID,omitempty"`
    Image_Url    url.URL `json:"Image_Url,omitempty" bson:"ID,omitempty"`
    Posted_timestamp time.Time  `json:"Posted_timestamp,omitempty" bson:"ID,omitempty"`
}

var client *mongo.client

func CreateUserEndpoint(response http.ResponseWriter,request *http.Request){
	response.Header().Add("Content-type","application/json")
	var user User
	json.NewDecoder(request,Body).Decode(&user)
	collection := client.database("insta").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx,user)
	json.NewEncoder(response).Encode(result)
}


func CreatePostEndpoint(response http.ResponseWriter,request *http.Request){
	response.Header().Add("Content-type","application/json")
	var post Post
	json.NewDecoder(request,Body).Decode(&user)
	collection := client.database("insta").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx,post)
	json.NewEncoder(response).Encode(result)
}


func GetPostsEndpoint(response http.ResponseWriter,request *http.request){
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _:= primitive.ObjectIDFromHex(params["id"])
	var post Post
    collection := client.Database("insta").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err := collection.FindOne(ctx,Post{ID: id}).Decode(&post)
	if err !=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

func GetUserPostsEndpoint(response http.ResponseWriter,request *http.request){
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _:= primitive.ObjectIDFromHex(params["id"])
	var userpost [] allpost
	collection := client.Database("Insta").Collection("userpost")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx,Post{ID: id}).Decode(&post)
	if err !=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(([]byte(`{"message":"`+err.Error()+`"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
        var post Post 
		cursor.Decode(&post)
		userpost = append(userpost,Post)
	}
	if err !=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(([]byte(`{"message":"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(response).Encode(userpost)
}

func GetUserEndpoint(response http.ResponseWriter,request *http.request){
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _:= primitive.ObjectIDFromHex(params["id"])
	var user User
    collection := client.Database("insta").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err := collection.FindOne(ctx,User{ID: id}).Decode(&user)
	if err !=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(([]byte(`{"message":"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func main(){
	fmt.println("Starting an application..")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.handleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.handleFunc("/User/{id}", GetUserEndpoint).Methods("GET")
	router.handleFunc("/post", CreatePostEndpoint).Methods("POST")
	router.handleFunc("/post/user/{id}", GetUserPostsEndpoint).Methods("GET")
	router.handleFunc("/Post/{id}", GetPostsEndpoint).Methods("GET")
	http.ListenAndServer(":12345",router)

}
