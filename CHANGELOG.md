[v0.1.5] 2024/10/28

Feature: AMQP Support; Lua Scripts Init; Analysis Module

- Added basic AMQP and AMQP exchange control
- Added Analysis Module Framework
- Added caller module of Persistent DB
- Added Direct XML Report

Update: Adapt bining simulation.go to conform to new lua script loading method

[v0.1.4] 2024/10/21
Update: Binding algorithm accuracy improvements
- Fix logic for achieving negative binding scores
- Removed deterministic vehicle locations in binding_simulation.go

[v0.1.3] 2024/10/21
Feature: Added binding_simulation.go
- Added binding simulation logic for binding analysis experiments

[v0.1.2] 2024/10/07
Feature: Added lua scripts for relationship binding
- Added binding lua scripts
- Updated binding go functions

[v0.1.1] 2024/09/30
Feature: MQTT, Redis, FCM, Gin clients and services
- Added MQTT client and handler functions
- Added Redis client, util functions and general lua scripts
- Added FCM service
- Added Gin webserver and serve logic
- Added module initializations in main.go

[v0.1.0] 2024/09/30
Initial Commit
- Added pkg folder and modules
- Added main.go