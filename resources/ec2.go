package resources


import (
    "log"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    
)


type EC2 struct {

  Instances *ec2.DescribeInstancesOutput
  Security_groups *ec2.DescribeSecurityGroupsOutput
  Volumes *ec2.DescribeVolumesOutput
  Snapshots *ec2.DescribeSnapshotsOutput
  last_err error

}


func (e *EC2) PurgeInstances(c1 chan int, meta interface{}, days int16, actualdestroy bool, whitelist []string) {

  sleepms := 800

  e.Instances, e.last_err = meta.(*AWSClient).ec2conn.DescribeInstances(nil)
  if e.last_err != nil {
      panic(e.last_err)
  }

  for idx, _ := range e.Instances.Reservations {
    for _, inst := range e.Instances.Reservations[idx].Instances {

      if !is_whitelisted(whitelist, *inst.InstanceId) {

        if(older_than(inst.LaunchTime.Unix(), days)) {

          if actualdestroy {

            params := &ec2.TerminateInstancesInput{
              InstanceIds: []*string{
                aws.String(*inst.InstanceId), 
              },
            }
            _, err := meta.(*AWSClient).ec2conn.TerminateInstances(params)

            if err != nil {
              log.Println("[ERROR] Instance: ", *inst.InstanceId, "could not be deleted, error occured: ", err )
            } else {
              log.Println("[INFO] Instance: ", *inst.InstanceId, " was terminated as it was older than ", days, " days.")
            }

            time.Sleep(time.Duration(sleepms) * time.Millisecond)

            sleepms += 600

          } else {
            log.Println("[WARN] Instance: ", *inst.InstanceId, " would be deleted without the dryrun!!!")
          }

        }
      }
    }
  }

  c1 <- 1

}


func (e *EC2) PurgeVolumes(c1 chan int, meta interface{}, days int16, actualdestroy bool, whitelist []string) {

  sleepms := 800

  params := &ec2.DescribeVolumesInput{}
  
  e.Volumes, e.last_err = meta.(*AWSClient).ec2conn.DescribeVolumes(params)

  if e.last_err != nil {
    panic(e.last_err)
  }

  for _, vl := range e.Volumes.Volumes {

    if *vl.State == "available" {

      if !is_whitelisted(whitelist, *vl.VolumeId) {

        if(older_than(vl.CreateTime.Unix(), days)) {

          if actualdestroy {

            params := &ec2.DeleteVolumeInput{
              VolumeId: aws.String(*vl.VolumeId), // Required
            }
            _, err := meta.(*AWSClient).ec2conn.DeleteVolume(params)

            if err != nil {
              log.Println("[ERROR] EBS Volume: ", *vl.VolumeId, "could not be deleted, error occured: ", err )
            } else {
              log.Println("[INFO] EBS Volume: ", *vl.VolumeId, " was deleted as it was older than ", days, " days.")
            }

            time.Sleep(time.Duration(sleepms) * time.Millisecond)

            sleepms += 600

          } else {
            log.Println("[WARN] EBS Volume: ", *vl.VolumeId, " would be deleted without the dryrun!!!")
          }

        }
      }
    }
  }

  c1 <- 1

}


func (e *EC2) PurgeSnapshots(c1 chan int, meta interface{}, days int16, actualdestroy bool, whitelist []string) {

  sleepms := 800

  params := &ec2.DescribeSnapshotsInput{
    OwnerIds: []*string{
      aws.String("self"),
    },
  }
  
  e.Snapshots, e.last_err = meta.(*AWSClient).ec2conn.DescribeSnapshots(params)

  if e.last_err != nil {
    panic(e.last_err)
  } 

  for _, sn := range e.Snapshots.Snapshots {

    if !is_whitelisted(whitelist, *sn.SnapshotId) {

      if(older_than(sn.StartTime.Unix(), days)) {

        if actualdestroy {

          params := &ec2.DeleteSnapshotInput{
            SnapshotId: aws.String(*sn.SnapshotId), // Required
          }
          _, err := meta.(*AWSClient).ec2conn.DeleteSnapshot(params)

          if err != nil {
            log.Println("[ERROR] Snapshot: ", *sn.SnapshotId, "could not be deleted, error occured: ", err )
          } else {
            log.Println("[INFO] Snapshot: ", *sn.SnapshotId, " was deleted as it was older than ", days, " days.")
          }

          time.Sleep(time.Duration(sleepms) * time.Millisecond)

          sleepms += 600

        } else {
          log.Println("[WARN] Snapshot: ", *sn.SnapshotId, " would be deleted without the dryrun!!!")
        }

      }
    }
  }

  c1 <- 1

}








