from typing import Optional

from auth import (
    session_token_generate,
    session_token_parts,
    session_token_verify,
)

admins = (
    (4, 'henrith', '$2b$12$tKBneY7fbLY4aMocM6OdOuhSoccx5o0dL/hJr.9ZqzraX1QetqF2O'),
    (10, 'seceng', '$2b$12$Cz1pcfw470GO7aUv/C/Ly.P0Ru7fGLcUBTj.1zSZdE0IwPy2ZSwz6'),
    (11, 'nick', '$2b$12$f8iLYhz0Vt4n0TKpG.bw2.MERr5uBlFShzzUztLZjb.TfQBmP2vLy'),
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
