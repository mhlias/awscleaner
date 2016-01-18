package main

import (
  "fmt"
  "flag"
  "time"

  "github.com/mhlias/awscleaner/resources"
)


func main() {
  
  awsconfigPtr := flag.String("config", "", "Location of aws config profiles to use")
  profilePtr := flag.String("profile", "", "Name of the profile to use")
  regionPtr := flag.String("region", "eu-west-1", "AWS Region to use")
  userolePtr := flag.Bool("userole", false, "Use instance role instead of aws config file")
  cfgfilePtr := flag.String("cfgfile", "config.yaml", "Location of custom config yaml file")
  destroyPtr := flag.Bool("destroy-resources", false, "If set to true all resources matching the config rules will be destroyed")



  flag.Parse()


  if ( (len(*awsconfigPtr) <= 0 || len(*profilePtr) <= 0) && !*userolePtr ) {
    fmt.Println("Please provide the following required parameters:")
    flag.PrintDefaults()
    return
  }


  cfg := &resources.Config{Awsconf: *awsconfigPtr, Profile: *profilePtr, Region: *regionPtr, Cfgfile: *cfgfilePtr }


  client := cfg.Connect()

  ec2 := new(resources.EC2)

  var fMap = map[string]func(c1 chan int, meta interface{}, days int16, actualrun bool, whitelist []string) {
    "ec2":  ec2.PurgeInstances,
    "ebs":  ec2.PurgeVolumes,
    "snapshots": ec2.PurgeSnapshots,
  }


  run := len(cfg.Conf.Resources)

  c1 := make(chan int)
  for k,v := range cfg.Conf.Resources {

    f, ok := fMap[k]

    if ok {
      go f(c1, client, v.Delete_older_than, *destroyPtr, v.Whitelist)
    }

  }

  
  completed := 0
  out := 0

  for {
    out = <- c1
    
    completed += out
    
    if(completed >= run ){
      break
    }

    time.Sleep(1)

  }


}