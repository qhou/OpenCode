package databases

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// A DocDB Object
type DocDB struct {
	Session        *mgo.Session
	IdeaCounter    int
	ProjectCounter int
}

// NewDocDB creates and returns a new DocumentDB connection
func NewDocDB(conn string) (DocDB, error) {
	// WORK AROUND FOR SSL //
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(conn)

	if err != nil {
		return DocDB{}, err
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	// Actual connection
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return DocDB{}, err
	}

	return DocDB{session, 0, 0}, nil
}

// ReturnProjects is...
func (d *DocDB) ReturnProjects() (*[]Project, error) {
	d.Session.Refresh()
	Coll := d.Session.DB("opencode").C("projects")
	var result []Project
	err := Coll.Find(nil).All(&result)
	return &result, err
}

// ReturnIdeas is...
func (d *DocDB) ReturnIdeas() (*[]Idea, error) {
	d.Session.Refresh()
	Coll := d.Session.DB("opencode").C("ideas")
	var result []Idea
	err := Coll.Find(nil).All(&result)
	return &result, err
}

// InsertProject is...
func (d *DocDB) InsertProject(Proj Project) error {
	d.Session.Refresh()
	Proj.ProjectID = d.ProjectCounter
	d.ProjectCounter++
	Proj.TimeStamp = time.Now()
	Coll := d.Session.DB("opencode").C("projects")
	err := Coll.Insert(&Proj)
	return err
}

// InsertIdea is...
func (d *DocDB) InsertIdea(Idea Idea) error {
	d.Session.Refresh()
	Idea.IdeaID = d.IdeaCounter
	d.IdeaCounter++
	Coll := d.Session.DB("opencode").C("ideas")
	err := Coll.Insert(&Idea)
	return err
}

// GetProjectByID is...
func (d *DocDB) GetProjectByID(id string) (*Project, error) {
	d.Session.Refresh()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	Coll := d.Session.DB("opencode").C("projects")
	var result Project
	err = Coll.Find(bson.M{"projectid": intID}).One(&result)
	return &result, err
}

// GetIdeaByID is...
func (d *DocDB) GetIdeaByID(id string) (*Idea, error) {
	d.Session.Refresh()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	Coll := d.Session.DB("opencode").C("ideas")
	var result Idea
	err = Coll.Find(bson.M{"ideaid": intID}).One(&result)
	return &result, err
}

// UpdateProjectEntry is ...
func (d *DocDB) UpdateProjectEntry(project *Project) error {
	d.Session.Refresh()
	Coll := d.Session.DB("opencode").C("projects")

	fmt.Println(project.ProjectID)
	colQuerier := bson.M{"projectid": project.ProjectID}
	change := bson.M{"$set": bson.M{"discussion": project.Discussion}}
	err := Coll.Update(colQuerier, change)
	return err
}

// UpdateIdeaEntry is ...
func (d *DocDB) UpdateIdeaEntry(idea *Idea) error {
	d.Session.Refresh()
	Coll := d.Session.DB("opencode").C("ideas")

	colQuerier := bson.M{"ideaid": idea.IdeaID}
	change := bson.M{"$set": bson.M{"discussion": idea.Discussion}}
	err := Coll.Update(colQuerier, change)
	return err
}
