import msgpack
import requests


message = {
    "Action": "SAVE",
    "Bucket": "webui",
    "Data": {
        "_id":"r2-test-01.png",
        "path": "/Users/harry/5.png",
        "mime": "image/png"
    }
}

m = msgpack.packb(message, use_bin_type=True)


sess = requests.session()
sess.keep_alive = False
url = "http://localhost:9090/r2/action"

r = sess.post(url, data=m, auth=('admin','123'), stream=True)
print(r.text)