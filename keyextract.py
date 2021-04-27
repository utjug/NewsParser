from gensim.summarization.summarizer import summarize
from flask import Flask, request

app = Flask(__name__)


@app.route('/', methods=['POST'])
def index():
    data = request.json
    sum = ""

    for x in data:
        res = summarize(x['content'], ratio=0.1)
        if res != "":
            print(res, end="\n\n")
            sum = sum + res + "\n\n"
        else:
            res = summarize(x['content'], ratio=0.5)
            sum = sum + res + "\n\n"

    return sum


app.run(host='0.0.0.0', port=5000)
