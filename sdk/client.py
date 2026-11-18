import json
import threading
import traceback
import urllib.error
import urllib.request
from dataclasses import dataclass


class PubSubError(Exception):
    def __init__(self, status, body):
        self.status = status
        self.body = body
        super().__init__(f"{status}: {body}")


@dataclass(frozen=True)
class Event:
    event: str
    message: str


class Topic:
    def __init__(self, client, name, register=True):
        self._client = client
        self.name = name
        self._handler = None
        self.listening = threading.Event()
        if register:
            client._register(name)

    def on(self, func):
        self._handler = func
        return func

    def publish(self, event, message):
        return self._client._publish(self.name, event, message)


class PubSubClient:
    def __init__(self, base_url):
        self._base_url = base_url.rstrip("/")
        self._topics = []

    def topic(self, name):
        created = Topic(self, name)
        self._topics.append(created)
        return created

    def subscribe(self, name):
        created = Topic(self, name, register=False)
        self._topics.append(created)
        return created

    def run(self):
        threads = []
        for topic in self._topics:
            if topic._handler is None:
                continue
            thread = threading.Thread(target=self._watch, args=(topic,), daemon=True)
            thread.start()
            threads.append(thread)
        for thread in threads:
            thread.join()

    def _register(self, name):
        return self._request("POST", "/register", {"name": name})

    def _publish(self, topic, event, message):
        return self._request(
            "POST",
            "/publish",
            {"topic": topic, "event": event, "message": message},
        )

    def _request(self, method, path, payload):
        response = self._open(method, path, payload)
        try:
            body = response.read().decode()
        finally:
            response.close()
        if not body:
            return ""
        data = json.loads(body)
        return data.get("message", "")

    def _open(self, method, path, payload):
        body = json.dumps(payload).encode()
        request = urllib.request.Request(
            self._base_url + path,
            data=body,
            headers={"Content-Type": "application/json"},
            method=method,
        )
        try:
            return urllib.request.urlopen(request)
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode().strip()
            raise PubSubError(exc.code, detail) from exc

    def _watch(self, topic):
        try:
            response = self._open("GET", "/subscribe", {"topic": topic.name})
        except PubSubError:
            traceback.print_exc()
            return
        topic.listening.set()
        try:
            for raw in response:
                line = raw.strip()
                if not line:
                    continue
                data = json.loads(line)
                event = Event(event=data.get("event", ""), message=data.get("message", ""))
                try:
                    topic._handler(event)
                except Exception:
                    traceback.print_exc()
        finally:
            response.close()
