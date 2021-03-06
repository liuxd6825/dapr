# -*- coding: utf-8 -*-
#
# Copyright 2021 The Dapr Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import asyncio

from flask import Flask, jsonify, request
from flask_dapr.actor import DaprActor

from dapr.conf import settings
from dapr.actor import ActorProxy, ActorId
from car_actor_interface import CarActorInterface
from car_actor import PythonCarActor

app = Flask(f'{PythonCarActor.__name__}Service')

# Enable DaprActor Flask extension
actor = DaprActor(app)
# Register DemoActor
actor.register_actor(PythonCarActor)

@app.route('/incrementAndGet/<actorType>/<actorId>', methods=['POST'])
def increment_and_get(actorType: str, actorId: str):
    print('/incrementAndGet/{0}/{1}'.format(actorType, actorId), flush=True)
    proxy = ActorProxy.create(actorType, ActorId(actorId), CarActorInterface)
    result = asyncio.run(proxy.IncrementAndGetAsync(1))
    return jsonify(result), 200

@app.route('/carFromJSON/<actorType>/<actorId>', methods=['POST'])
def car_from_json(actorType: str, actorId: str):
    print('/carFromJSON/{0}/{1}'.format(actorType, actorId), flush=True)
    proxy = ActorProxy.create(actorType, ActorId(actorId), CarActorInterface)
    body = request.get_data().decode('utf-8')
    result = asyncio.run(proxy.CarFromJSONAsync(body))
    return jsonify(result), 200

@app.route('/carToJSON/<actorType>/<actorId>', methods=['POST'])
def car_to_json(actorType: str, actorId: str):
    print('/carToJSON/{0}/{1}'.format(actorType, actorId), flush=True)
    proxy = ActorProxy.create(actorType, ActorId(actorId), CarActorInterface)
    car = request.get_json()
    result = asyncio.run(proxy.CarToJSONAsync(car))
    return result, 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=settings.HTTP_APP_PORT)
