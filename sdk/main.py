from client import PubSubClient


def handle_order(event):
    print(event.event, event.message)


def main():
    client = PubSubClient("http://localhost:8080")
    orders = client.subscribe("orders")
    orders.on(handle_order)
    client.run()


if __name__ == "__main__":
    main()
