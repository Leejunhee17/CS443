import argparse
import json

import requests


class CLI:
    def __init__(self):
        self.local_addr = "http://127.0.0.1:5001/internal"

    def ask_account(self, data):
        return requests.post(self.local_addr + "/account/", json=data)

    def ask_check_balance(self):
        return requests.get(self.local_addr + "/balance/")

    def ask_publish_tx(self, data):
        return requests.post(self.local_addr + "/tx/", json=data)

    def ask_start_mining(self):
        return requests.post(self.local_addr + "/mining/start/")

    def ask_stop_mining(self):
        return requests.post(self.local_addr + "/mining/stop/")

    @staticmethod
    def response_handler(resp):
        code = resp.status_code
        print(resp)
        if code == 200:
            body = resp.json()
            if body is not None:
                print(body)
            return True
        elif code == 403:
            body = resp.json()
            print(code, body["msg"])
        try:
            resp.raise_for_status()
        except requests.exceptions.HTTPError as e:
            print(e)
        return False


if __name__ == '__main__':
    cli = CLI()
    print("Set account.\n")
    while True:
        fname = input("File name to create[load] an account: ").strip()
        resp = cli.ask_account({"fname": fname})
        if cli.response_handler(resp):
            break

    while True:
        op = int(
            input("Enter number.\n" + "\t1. Check balance\n" +
                  "\t2. Send money\n" + "\t3. Start mining\n" +
                  "\t4. Stop mining\n" + "\t0. Quit\n"))
        if op == 0:
            break
        elif op == 1:
            resp = cli.ask_check_balance()
        elif op == 2:
            resp = cli.ask_publish_tx({
                "receiver":
                input("Receiver: \n"),
                "amount":
                int(input("Amount: \n")),
            })
        elif op == 3:
            resp = cli.ask_start_mining()
        elif op == 4:
            resp = cli.ask_stop_mining()
        else:
            print("Invalid input.\n" + "Please enter number between 1 and 5.")
            continue
        cli.response_handler(resp)
