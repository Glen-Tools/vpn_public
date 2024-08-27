#!/bin/bash

docker exec wg-dante sh -c "wg-quick down wg0; killall danted"