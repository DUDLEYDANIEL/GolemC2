# Local Golem C2

This project is a locally deployed C2 server like framework built in Go. It simulates a centralized management system that coordinates multiple remote agents in a controlled lab environment. The goal is to study secure communication, authentication models, distributed task orchestration, and adversary-emulation concepts from a defensive engineering perspective.

The system consists of three main components:
* A central server that manages agents, assigns tasks, and collects results
* Lightweight agents written in Go that periodically connect to the server, retrieve tasks, execute them, and return results
* A CLI tool that allows operators to issue commands and monitor activity

All communication between server and agents is designed to occur over HTTPS with mutual TLS (mTLS), ensuring encrypted transport and strong identity verification. The architecture also supports per-agent asymmetric key pairs (RSA or Ed25519) for enhanced authentication and secure result handling. Key management, certificate validation, and transport-layer security are core design considerations.

The framework supports controlled lab simulations of:
* Host reconnaissance
* Network scanning and service identification
* Internal network mapping
* Secure result transmission

This project focuses on:
* Secure distributed system design
* Certificate-based authentication
* Cryptographic key management
* gRPC-based service architecture (in development)

Defensive research and detection-aware engineering

It is strictly intended for educational and research purposes within isolated lab environments. The system is not designed for unauthorized deployment or real-world misuse.

## Installation

Install my-project with npm

```bash
git clone https://github.com/DUDLEYDANIEL/GolemC2.git
cd GolemC2
```

## Usage/Examples

```bash
curl --cert certs/client-cert.pem \
     --key certs/client-key.pem \
     --cacert certs/ca-cert.pem -X POST "https://localhost:8443/register?agent_id=test2-agent"

curl --cert certs/client-cert.pem \
     --key certs/client-key.pem \
     --cacert certs/ca-cert.pem \
  -X POST "https://localhost:8443/tasks?agent_id=test2-agent" \
  -H "Content-Type: application/json" \
  -d '{"type":"execute","params":{"command":"whoami && pwd && ls -la"}}'

curl --cert certs/client-cert.pem \
     --key certs/client-key.pem \
     --cacert certs/ca-cert.pem "https://localhost:8443/tasks?agent_id=test2-agent"

..........
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
