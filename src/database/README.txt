func Connect_localhost(URI string, user string, password string) *mongo.Client -> Connect to localhost using uri, username and password, returns client

Connect_cluster(user string, password string, database string) -> Connect to cluster using username, password and database name, returns client

Disconnect(client *mongo.Client) -> Disconnect client

Insert_to_collection(client *mongo.Client, database string, col string, path string) bool -> Insert json file form path to a database collection. The collection with name "col" is created if it didn't exist

Delete_documents(client *mongo.Client, database string, col string, filter string, value string) bool -> Delete documents that match the value "value" of the filter "filter" from collection "col"

Update_documents(client *mongo.Client, database string, col string, filter string, filter_value string,update_Param string, update_value string) bool -> Update field "update_Param" with the value "update_value" from docuements that match the value "value" of the filter "filter" from collection "col" 


