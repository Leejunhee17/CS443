from flask import Flask, Response, request

from node.node import Node

from . import PORT

node = Node()
app = Flask(__name__)


@app.route('/internal/account/', methods=["POST"])
def post_internal_account():
    print("post internal account")
    data = request.get_json()
    (res, status_code) = node.post_internal_account(data)
    return Response(res, status_code)


@app.route('/internal/balance/', methods=['GET'])
def get_internal_balance():
    print("get internal balance")
    (res, status_code) = node.get_internal_balance()
    return Response(res, status_code)


@app.route('/internal/tx/', methods=["POST"])
def post_internal_tx():
    print("post internal tx")
    data = request.get_json()
    (res, status_code) = node.post_internal_tx(data)
    return Response(res, status_code)


@app.route('/internal/mining/start/', methods=["POST"])
def post_internal_mining_start():
    print("post internal mining start")
    (res, status_code) = node.post_internal_mining_start()
    return Response(res, status_code)


@app.route('/internal/mining/stop/', methods=["POST"])
def post_internal_stop_mining():
    print("post stop mining")
    (res, status_code) = node.post_internal_mining_stop()
    return Response(res, status_code)


@app.route('/peer/block', methods=["POST"])
def post_peer_block():
    print("post peer block")
    data = request.get_json()
    (res, status_code) = node.post_peer_block(data)
    return Response(res, status_code)


@app.route('/peer/tx', methods=["POST"])
def post_peer_tx():
    print("post peer tx")
    data = request.get_json()
    (res, status_code) = node.post_peer_tx(data)
    return Response(res, status_code)


def manager_run():
    app.run(host='0.0.0.0', debug=True, port=PORT)
