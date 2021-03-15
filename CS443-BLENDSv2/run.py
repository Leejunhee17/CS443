from blends.node.node import Node

node = Node(
    "http://143.248.47.200:8080/",
    "my_key.json",
    "sqlite:///scenario.db",
)

node.run()