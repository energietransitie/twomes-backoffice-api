from typing import Optional

from auth import (
    session_token_generate,
    session_token_parts,
    session_token_verify,
)

admins = (
    (1, 'arjan', '$2b$12$BAkuzJVXNwgeiU2vBL5b0eBPDm777eI0bvnjqKmlpDNGBr/hdGnjq'),
    (2, 'Marco', '$2b$12$mP7XyiffPKDGJM05xGPrBuaE8yoLVjabaXj1.UYaxFYdd3bpLzUwu'),
    (3, 'Kevin', '$2b$12$9iG1PXmGN9BJHyypW6rMgOT2wHx/zrJhg/rZT0/GgxTnwimnl1CBC'),
    (4, 'henrith', '$2b$12$tKBneY7fbLY4aMocM6OdOuhSoccx5o0dL/hJr.9ZqzraX1QetqF2O'),
    (5, 'Tristan', '$2b$12$DdmloxiPuI3YDBM5YxmksO7JShR1KR8qc9WFUbLBVcPhAq/tNIeNO'),
    (6, 'apptest', '$2b$12$WW3OcYnfqE.Xc6d1v/8sceduMPr3LENkKIBPWR46JfJ9sa8q6nXk2')
    # Add extra admin tuples, created with 'create_admin()', here
)
admins_name = {a[0]: a[1] for a in admins}
admins_session_token_hash = {a[0]: a[2] for a in admins}


def create_admin():
    admin_id = max([a[0] for a in admins]) + 1 if admins else 1
    admin_name = input('Admin name: ')

    admin_token, admin_token_hash = session_token_generate(admin_id)

    admin = (admin_id, admin_name, admin_token_hash)

    print(f'Update {__name__}.admins with this tuple: \n  {admin}')
    print(f'Authorisation bearer token for admin "{admin_name}": \n  {admin_token}')


def get_admin(admin_session_token: str) -> Optional[str]:
    try:
        admin_id, _ = session_token_parts(admin_session_token)
        admin_name = admins_name[admin_id]
        admin_token_hash = admins_session_token_hash[admin_id]

    except Exception:
        return None

    if not session_token_verify(admin_session_token, admin_token_hash):
        return None

    return admin_name
