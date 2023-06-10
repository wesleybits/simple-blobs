import http
import http.server
from types import TracebackType
from typing_extensions import override
import urllib
import urllib.parse
import typing
import abc
import json

# Let's do some prototyping!! We're gonna be doing what looks like enterprise
# anti-patterns, _but_ it's completely to demonstrate using some intermediary
# "type" to abstract the actual working code, so that the working code can be
# readily replaced without needing to alter the entire project.

T = typing.TypeVar('T')
class RepoInterface(typing.Generic[T], metaclass=abc.ABCMeta):
    @classmethod
    def __subclasshook__(cls, __subclass) -> bool:
        return (hasattr(__subclass, "get") and callable(__subclass.get) and
                hasattr(__subclass, "put") and callable(__subclass.put) and
                hasattr(__subclass, "delete") and callable(__subclass.delete) and
                hasattr(__subclass, "values") and callable(__subclass.values) or
                NotImplemented)

    @abc.abstractmethod
    def get(self, id: str) -> T:
        raise NotImplementedError

    @abc.abstractmethod
    def put(self, id: str, data: T) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def delete(self, id: str) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def values(self) -> list[T]:
        raise NotImplementedError

class ModelInterface(typing.Generic[T], metaclass=abc.ABCMeta):
    @classmethod
    def __subclasshook__(cls, __subclass) -> bool:
        return (hasattr(__subclass, 'create') and callable(__subclass.create) and
                hasattr(__subclass, 'update') and callable(__subclass.update) and
                hasattr(__subclass, 'delete') and callable(__subclass.delete) and
                hasattr(__subclass, 'read_one') and callable(__subclass.read_one) and
                hasattr(__subclass, 'select') and callable(__subclass.select) or
                NotImplemented)

    @abc.abstractmethod
    def create(self, data: T) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def update(self, id: str, data: T) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def delete(self, id: str) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def read_one(self, id: str) -> T:
        raise NotImplementedError

    @abc.abstractmethod
    def select(self, query: dict) -> T:
        raise NotImplementedError

class Controller(metaclass=abc.ABCMeta):
    @classmethod
    def __subclasshook__(cls, __subclass) -> bool:
        return (hasattr(__subclass, "search") and callable(__subclass.search) and
                hasattr(__subclass, "get_one") and callable(__subclass.get_one) and
                hasattr(__subclass, "post") and callable(__subclass.post) and
                hasattr(__subclass, "put") and callable(__subclass.put) and
                hasattr(__subclass, "delete") and callable(__subclass.delete) and
                hasattr(__subclass, "root_handler") and callable(__subclass.root_handler) or
                NotImplemented)

    @abc.abstractmethod
    def root_handler(self, root_handler: http.server.BaseHTTPRequestHandler) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def search(self, parameters: dict) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def get_one(self, id: str) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def post(self) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def put(self, id: str) -> None:
        raise NotImplementedError

    @abc.abstractmethod
    def delete(self, id: str) -> None:
        raise NotImplementedError

# Prototypes done, now let's create our dispatcher. Again, it's only going
# to _depend_ on the above interfaces. Not only does it make it uncomplicated
# to build mocks, but it also ensures that I can write this without actually
# implementing anything at all and pyright is 100% A-OK with it!

class HTTPNotFound(Exception):
    pass

class Handler(http.server.BaseHTTPRequestHandler):
    __prefixes: typing.ClassVar[list[tuple[str, typing.Type[Controller]]]] = []

    @classmethod
    def add_controller(cls, prefix: str, controller: typing.Type[Controller]) -> None:
        cls.__prefixes += [(prefix, controller)]

    def find_controller(self) -> tuple[str, Controller, str, dict[str, str]]:
        parsed = urllib.parse.urlparse(self.path)
        for prefix, controlcls in self.__prefixes:
            if self.path.startswith(parsed[2]):
                ctrl = controlcls()
                id = parsed[2].lstrip(prefix)
                params = dict([p.split("=") for p in parsed[4].split('&')])
                ctrl.root_handler(self)
                return prefix, ctrl, id, params
        else:
            raise HTTPNotFound

    @override
    def handle_one_request(self) -> None:
        try:
            return super().handle_one_request()
        except HTTPNotFound:
            self.send_response(404)
        except NotImplementedError:
            self.send_response(405)
        except Exception as e:
            payload_bytes = e.with_traceback(None).__str__().encode()
            self.send_response(500)
            self.send_header('Content-Type', 'plain/text')
            self.send_header('Content-Length', str(len(payload_bytes)))
            self.end_headers()
            self.wfile.write(payload_bytes)

    def reply_ok(self, payload: typing.Union[list, dict, str]) -> None:
        payload_bytes = json.dumps(payload).encode()
        payload_length = len(payload_bytes)
        self.send_response(200)
        self.send_header('Content-Type', 'plain/text')
        self.send_header('Content-Length', str(payload_length))
        self.end_headers()
        self.wfile.write(payload_bytes)

    def do_GET(self) -> None:
        _, ctrl, id, q = self.find_controller()
        if not id:
            return ctrl.search(q)
        else:
            return ctrl.get_one(id)

    def do_PUT(self) -> None:
        _, ctrl, id, _ = self.find_controller()
        if id:
            return ctrl.put(id)
        else:
            raise HTTPNotFound

    def do_POST(self) -> None:
        _, ctrl, _, _ = self.find_controller()
        return ctrl.post()

    def do_DELETE(self) -> None:
        _, ctrl, id, _ = self.find_controller()
        if id:
            return ctrl.delete(id)
        else:
            raise HTTPNotFound

    def do_ORIGIN(self) -> None:
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', '*')
        self.send_header('Access-Control-Allow-Headers', '*')
        self.end_headers()

