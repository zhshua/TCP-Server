package main

import (
	"TCP-Server/protobufDemo/pb"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "zhangshaui",
		Age:    18,
		Emails: []string{"zs1@163.com", "zs2@163.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "111111111",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "222222222",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "3333333333",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	// 将protobuf的message进行序列化, 得到一个二进制文件
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("proto.Marshal err", err)
	}

	newperson := &pb.Person{}
	// 将data反序列化
	err = proto.Unmarshal(data, newperson)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
	}
	fmt.Println("源数据：", person)
	fmt.Println("解码之后的数据：", newperson)

}
