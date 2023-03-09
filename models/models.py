from typing import List

class MemberConfig(dict):
    username: str
    has_joined_voice: bool
    has_joined_text: bool
    # MAYBE INSERT OTHER DATETIME DELTA PROPERTIES HERE
    active: bool
