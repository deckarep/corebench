package providers

// import (
//   "./aws_templates.go"
// )



const (
AwsProviderInstanceNameFmt = "corebench-aws-%s"
AwsBenchReadyScript     = "export GOPATH=/home/ubuntu/go && while [ ! -f $GOPATH/.core-init ]; do sleep 1; done"
AwsBenchCommandTemplate = `cd $GOPATH/src/${git-repo} && /usr/local/go/bin/go get . && /usr/local/go/bin/go version && /usr/local/go/bin/go test -v ${benchmem-setting}-cpu ${cpu-count} -bench=${bench-regex} -count=${bench-count}`
AwsBenchStatTemplate    = " | tee benchmark.log && echo '\n\n' && $GOPATH/bin/benchstat benchmark.log"
)

var (
	awsGoVersionFmt      = "go%s.linux-amd64.tar.gz"
)

const CfnTemplate = `
---
AWSTemplateFormatVersion: '2010-09-09'
Description: ''
Parameters:
  KeyName:
    Type: AWS::EC2::KeyPair::KeyName
    Default: ${keypair}
  ImageId:
    Type: String
    Default: ami-d944fda6
  #  Default: ami-d38a4ab1
  InstanceType:
    Type: String
    Default: ${instancetype}

Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: 172.17.0.0/16
      InstanceTenancy: default
      EnableDnsSupport: true
      EnableDnsHostnames: true

  Subnet1:
    Type: AWS::EC2::Subnet
    Properties:
      CidrBlock: 172.17.1.0/24
      AvailabilityZone: ${awsregion}a
      VpcId: !Ref VPC

  SecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: corebench
      VpcId: !Ref VPC

  InboundSGEntry1:
    Type: AWS::EC2::SecurityGroupIngress
    Properties:
      GroupId: !Ref SecurityGroup
      IpProtocol: tcp
      FromPort: 22
      ToPort: 22
      CidrIp: 0.0.0.0/0

  InboundSGEntry2:
    Type: AWS::EC2::SecurityGroupIngress
    Properties:
      GroupId: !Ref SecurityGroup
      IpProtocol: tcp
      FromPort: 22
      ToPort: 22
      CidrIp: 172.17.1.0/24

  InternetGateway:
    Type: AWS::EC2::InternetGateway
    Properties:
      Tags:
      - Key: Name
        Value: whatever

  VPCGatewayAttachment:
    Type: AWS::EC2::VPCGatewayAttachment
    Properties:
      VpcId: !Ref VPC
      InternetGatewayId: !Ref InternetGateway

  Route1:
    Type: AWS::EC2::Route
    Properties:
      DestinationCidrBlock: 0.0.0.0/0
      RouteTableId: !Ref RouteTable
      GatewayId: !Ref InternetGateway

  NetworkAcl:
    Type: AWS::EC2::NetworkAcl
    Properties:
      VpcId: !Ref VPC

  ACL1:
    Type: AWS::EC2::NetworkAclEntry
    Properties:
      CidrBlock: 0.0.0.0/0
      Egress: 'true'
      Protocol: "-1"
      RuleAction: allow
      RuleNumber: '100'
      NetworkAclId: !Ref NetworkAcl

  ACL2:
    Type: AWS::EC2::NetworkAclEntry
    Properties:
      CidrBlock: 0.0.0.0/0
      Protocol: "-1"
      RuleAction: allow
      RuleNumber: '100'
      NetworkAclId: !Ref NetworkAcl

  SubnetACL1:
    Type: AWS::EC2::SubnetNetworkAclAssociation
    Properties:
      NetworkAclId: !Ref NetworkAcl
      SubnetId: !Ref Subnet1

  RouteTable:
    Type: AWS::EC2::RouteTable
    Properties:
      VpcId: !Ref VPC

  SubnetRoute1:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      RouteTableId: !Ref RouteTable
      SubnetId: !Ref Subnet1

  corebench:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: !Ref ImageId
      InstanceType: !Ref InstanceType
      KeyName: !Ref KeyName
      NetworkInterfaces:
      - SubnetId: !Ref Subnet1
        AssociatePublicIpAddress: 'true'
        DeviceIndex: '0'
        GroupSet:
        - !Ref SecurityGroup
      UserData:
          Fn::Base64: !Sub |
            #!bin/bash -xe
            apt-get update

            echo "Setting up aws tools for the first time..."
            apt-get -y install gv awscli zip python-setuptools
            mkdir aws-cfn-bootstrap-latest
            curl https://s3.amazonaws.com/cloudformation-examples/aws-cfn-bootstrap-latest.tar.gz | tar xz -C aws-cfn-bootstrap-latest --strip-components 1
            easy_install aws-cfn-bootstrap-latest

            echo "Setting up corebench for the first time..."
            echo "Installing dependencies..."
            apt-get -y install git
            wget https://storage.googleapis.com/golang/${go-version}
            tar -C /usr/local -xzf ${go-version}
            export GOROOT=/usr/local/go
            export GOPATH=/home/ubuntu/go
            mkdir -p $GOPATH
            $GOROOT/bin/go get github.com/golang/perf/cmd/benchstat
            $GOROOT/bin/go get ${git-repo}
            touch $GOPATH/.core-init
            chown -R ubuntu:ubuntu $GOPATH
            echo "Finished corebench initialization"

            if [ $? != 1 ]; then
              echo "Signalling stack complete"
              state=0
              /usr/local/bin/cfn-signal -e $state --stack ${AWS::StackName} --resource corebench --region ${AWS::Region}
            else
              echo "Signalling stack create failed"
              state=1
              /usr/local/bin/cfn-signal -e $state --stack ${AWS::StackName} --resource corebench --region ${AWS::Region}
            fi
      Tags:
        -
          Key: role
          Value: corebench
    CreationPolicy:
      ResourceSignal:
        Count: 1
        Timeout: PT5M

`
