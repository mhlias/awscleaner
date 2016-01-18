package resources


import (
    "fmt"
    "log"
    "io/ioutil"
    "path/filepath"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/credentials"

    "gopkg.in/yaml.v2"
)



type AWSClient struct {
  ec2conn            *ec2.EC2

  region             string
}


type Config struct {
  Awsconf string
  Profile string
  Region  string
  Cfgfile string
  Call  map[string] func()
  Conf    yamlConfig
}


type yamlConfig struct {
  Resources map[string] struct {
    Delete_older_than int16
    Whitelist []string
  }
}



func (c *Config) Connect() interface{} {


  c.readConf()

  var client AWSClient
  
  awsConfig := new(aws.Config)

  if len(c.Profile)>0 {
    awsConfig = &aws.Config{
      Credentials: credentials.NewSharedCredentials(c.Awsconf, fmt.Sprintf("profile %s", c.Profile)),
      Region:      aws.String(c.Region),
      MaxRetries:  aws.Int(3),
    }

  } else {
    // use instance role
    awsConfig = &aws.Config{
      Region:      aws.String(c.Region),
    }

  }


  sess := session.New(awsConfig)

  client.ec2conn = ec2.New(sess)


  return &client

}

func (c *Config) readConf() {

  configFile, _ := filepath.Abs(c.Cfgfile)
  yamlConf, file_err := ioutil.ReadFile(configFile)

  if (file_err != nil) {
    log.Println("[ERROR] File does not exist: ", file_err)
  }

  yaml_err := yaml.Unmarshal(yamlConf, &c.Conf)

  if (yaml_err != nil) {
    panic(yaml_err)
  }

}


