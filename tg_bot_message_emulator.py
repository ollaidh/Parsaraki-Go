import requests
from datetime import datetime
import time


def send_message_to_tg_handler(url: str):
    curr_text = str(datetime.now())
    payload = {
        "update_id": 859056273,
        "message": {
            "message_id": 143,
            "from": {
                "id": 1361646730,
                "is_bot": False,
                "first_name": "Maria",
                "last_name": "Lineva",
                "username": "Maria_Lineva",
                "language_code": "ru",
                "is_premium": True,
            },
            "chat": {
                "id": 1361646730,
                "first_name": "Maria",
                "last_name": "Lineva",
                "username": "Maria_Lineva",
                "type": "private",
            },
            "date": 1771529141,
            "text": curr_text,
        },
    }
    headers = {
        "X-Telegram-Bot-Api-Secret-Token": "fwKklvjtDqHHP44KwZOMtoJXvbbf7tVt6fxJqjfwJ5mZBBvXqlv2rR19G5WdFT3w",
    }
    result = requests.post(url=url, json=payload, headers=headers)
    print(result.status_code, curr_text)


if __name__ == "__main__":
    url = "http://127.0.0.1:8443/bot-message"
    for i in range(5):
        send_message_to_tg_handler(url)
        time.sleep(1)
