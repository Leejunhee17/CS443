import threading
import queue
from urllib.parse import urljoin

import requests
import time


class NetworkManager(threading.Thread):
    def __init__(self, broadcast_url: str, data_queue: queue.Queue):
        threading.Thread.__init__(self)
        self.broadcast_url = broadcast_url
        self.data_queue = data_queue
        self.last_time = 0

    def fetch(self):
        results = requests.get(
            urljoin(self.broadcast_url, "/get"),
            params={
                "start_at": self.last_time
            }).json()

        self.last_time += len(results)
        return results

    def publish(self, data):
        return requests.post(urljoin(self.broadcast_url, "/post"), json=data)

    def run(self):
        while True:
            results = self.fetch()
            self.data_queue.put(results)
            if len(results) > 0:
                continue
            time.sleep(60)