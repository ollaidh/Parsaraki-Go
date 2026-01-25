from dotenv import load_dotenv, set_key
from pyngrok import ngrok
import time


def start_ngrok() -> str:
    tunnels = ngrok.get_tunnels()
    if len(tunnels) > 0:
        result = tunnels[0].public_url
        print("Using existing tunnel:", result)
        return result

    ngrok.connect(8443)
    tunnels = ngrok.get_tunnels()
    result = tunnels[0].public_url

    print("Created new tunnel, public URL:", result)
    print("Local URL:", tunnels[0].config["addr"])
    return result


def update_config_with_url(url: str, env_filepath: str):
    load_dotenv(env_filepath)
    set_key(env_filepath, "WEBHOOKS_URL", url)


if __name__ == "__main__":
    ngrok_public_url = start_ngrok()
    update_config_with_url(ngrok_public_url, ".env")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        ngrok.kill()
        print("\nNgrok stopped!")
