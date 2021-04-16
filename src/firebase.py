from urllib.parse import urlencode, urlunsplit

# Placeholder values, to construct the Firebase Dynamic Link
# TODO: update with the final values
FIREBASE_APP = 'energietransitiewindesheim'
FIREBASE_HOST = f'{FIREBASE_APP}.page.link/'

FIREBASE_ANDROID_PACKAGE = 'nl.windesheim.androidapp'
FIREBASE_IOS_BUNDLE = 'nl.windesheim.iosapp'

API_URL = 'https://api.tst.energietransitiewindesheim.nl'
APP_URL = 'etw://account/{activation_token}'


def firebase_dynamic_link(activation_token: str) -> str:
    """
    Generate a Firebase Dynamic Link, using an account activation token.
    """
    query_params = {
        'link': APP_URL.format(activation_token=activation_token),
        'apn': FIREBASE_ANDROID_PACKAGE,
        'isi': FIREBASE_IOS_BUNDLE,
    }
    query_string = urlencode(query_params)

    url = urlunsplit(['https', FIREBASE_HOST, '', query_string, ''])
    return url
