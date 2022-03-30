from urllib.parse import urlencode, urlunsplit

FIREBASE_APP = 'energietransitiewindesheim'
FIREBASE_HOST = f'{FIREBASE_APP}.page.link/'

FIREBASE_ANDROID_PACKAGE = 'nl.windesheim.energietransitie.warmtewachter'
FIREBASE_IOS_BUNDLE = 'nl.windesheim.energietransitie.warmtewachter'
FIREBASE_IOS_STORE_ID = '1563201993'
FIREBASE_IOS_SKIP_PREVIEW = '1'

API_URL = 'https://api.tst.energybehaviour.net'
APP_URL = 'https://account/{activation_token}'


def firebase_dynamic_link(activation_token: str) -> str:
    """
    Generate a Firebase Dynamic Link, using an account activation token.
    """
    query_params = {
        'link': APP_URL.format(activation_token=activation_token),
        'apn': FIREBASE_ANDROID_PACKAGE,
        'ibi': FIREBASE_IOS_BUNDLE,
        'isi': FIREBASE_IOS_STORE_ID,
        'efr': FIREBASE_IOS_SKIP_PREVIEW
    }
    query_string = urlencode(query_params)

    url = urlunsplit(['https', FIREBASE_HOST, '', query_string, ''])
    return url
