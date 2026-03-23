import copy
import yaml
import sys

COMPOSE_CONTENT_TEMPLATE = {
    "name": "tp0",
    "services": {
        "server": {
            "container_name": "server",
            "image": "server:latest",
            "entrypoint": "python3 /main.py",
            "environment": ["PYTHONUNBUFFERED=1"],
            "networks": ["bet_network"],
            "volumes": ["./server/config.ini:/config.ini"],
        }
    },
    "networks": {
        "bet_network": {
            "ipam": {
                "driver": "default",
                "config": [{"subnet": "172.25.125.0/24"}],
            }
        }
    },
}

CLIENT_TEMPLATE = {
    "image": "client:latest",
    "entrypoint": "/client",
    "networks": ["bet_network"],
    "depends_on": ["server"],
    "volumes": ["./client/config.yaml:/config.yaml"],
}


def generate_compose_file(file_name: str, client_count: int) -> None:
    compose_content = copy.deepcopy(COMPOSE_CONTENT_TEMPLATE)

    for i in range(1, client_count + 1):
        service_name = f"client{i}"
        client_content = copy.deepcopy(CLIENT_TEMPLATE)
        client_content["container_name"] = service_name
        client_content["environment"] = [f"CLI_ID={i}"]
        compose_content["services"][service_name] = client_content

    with open(file_name, "w") as file:
        yaml.dump(compose_content, file)


def main() -> None:
    if len(sys.argv) != 3:
        print("Usage: python generate-compose.py <output_file_name> <client_count>")
        sys.exit(1)

    output_file_name = sys.argv[1]
    client_count = int(sys.argv[2])
    generate_compose_file(output_file_name, client_count)


if __name__ == "__main__":
    main()
