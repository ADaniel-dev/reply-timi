"""
请确保在运行代码之前先根据给定的Protobuf文件定义生成Python代码（使用protoc命令或相关插件）。然后，将生成的addressbook_pb2.py文件与上述代码放在同一个目录中。
"""

import addressbook_pb2

# 创建一条新的Person记录
person = addressbook_pb2.Person()
person.name = "John Doe"
person.id = 123
person.email = "john.doe@example.com"

# 添加电话号码
phone = person.phones.add()
phone.number = "1234567890"
phone.type = addressbook_pb2.Person.HOME

# 创建AddressBook并添加Person记录
address_book = addressbook_pb2.AddressBook()
address_book.people.append(person)

# 写入二进制文件
with open("addressbook.bin", "wb") as f:
    f.write(address_book.SerializeToString())

# 从二进制文件中读取数据
with open("addressbook.bin", "rb") as f:
    data = f.read()

# 解析二进制数据
address_book = addressbook_pb2.AddressBook()
address_book.ParseFromString(data)

# 遍历Person记录并显示内容
for person in address_book.people:
    print("Name:", person.name)
    print("ID:", person.id)
    print("Email:", person.email)
    print("Phone Numbers:")
    for phone in person.phones:
        print("- Number:", phone.number)
        print("  Type:", addressbook_pb2.Person.PhoneType.Name(phone.type))
    print("--------------------")
