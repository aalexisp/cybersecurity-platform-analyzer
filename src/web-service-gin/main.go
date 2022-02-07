package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type box_specifications struct {
	CPU_Number	int	`json:"cpu_number"`
	Arch	string	`json:"architecture"`
	Ram		string	`json:"ram"`
	Provider	string	`json:"provider"`
	TTL		string	`json:"ttl"`
	Box_ip	string	`json:"box_ip"`
	Gui		bool	`json:"bool"`
}

// album represents data about a record album.
type parameters struct {
	ID     string  `json:"id"`
	Name  string  `json:"name"`
	Host_ip string  `json:"host_ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Box_spec box_specifications `json:"box_spec"`
}

var machines = []parameters{
	{ID: "1", Name: "Slow", Host_ip:"8.8.8.8.8", Username:"root",
	Password:"toor", Box_spec: box_specifications{CPU_Number:1, Arch:"x86_64", Ram:"256",
	Provider:"??", TTL:"9999", Box_ip:"8.8.8.8.8", Gui:true}},

	{ID: "2", Name: "Medium", Host_ip:"8.8.8.8.8", Username:"root",
	Password:"toor", Box_spec: box_specifications{CPU_Number:2, Arch:"x86_64", Ram:"512",
	Provider:"??", TTL:"9999", Box_ip:"8.8.8.8.8", Gui:true}},

	{ID: "3", Name: "Fast", Host_ip:"8.8.8.8.8", Username:"root",
	Password:"toor", Box_spec: box_specifications{CPU_Number:4, Arch:"x86_64", Ram:"1024",
	Provider:"??", TTL:"9999", Box_ip:"8.8.8.8.8", Gui:true}},
}

func getMachines(c *gin.Context){
	c.IndentedJSON(http.StatusOK, machines)
}

func postMachines(c *gin.Context) {
	var newMachine parameters

	// Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newMachine); err != nil {
        return
    }

    // Add the new album to the slice.
    machines = append(machines, newMachine)
    c.IndentedJSON(http.StatusCreated, newMachine)
}

func main() {
	router := gin.Default()
	router.GET("/machines", getMachines)
	router.POST("/machines", postMachines)

	router.Run("localhost:8080")
}
