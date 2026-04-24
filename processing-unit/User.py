import json
import os
import time
from typing import Any


def _load_settings() -> dict[str, float]:
    settings_path = os.path.join(os.path.dirname(__file__), "setting.json")
    try:
        with open(settings_path, "r", encoding="utf-8") as settings_file:
            return json.load(settings_file)
    except FileNotFoundError:
        return {}


settings = _load_settings()
AUTH_COOLDOWN = settings.get("AUTH_COOLDOWN", 60.0)
GESTURE_COOLDOWN = settings.get("GESTURE_COOLDOWN", 3.0)


class User:
    def __init__(self, user_id: str, embedding: Any):
        self.user_id = user_id
        self.embedding = embedding
        self.auth_token: str | None = None

        # Timestamp when auth was last successfully issued.
        self.last_auth_time: float | None = None

        # Timestamp when an auth attempt was last made (success or failure).
        self.last_auth_attempt_time: float | None = None

        # Gesture ID -> timestamp of last successful workflow trigger.
        self.gesture_timestamps: dict[int, float] = {}

        # Backward compatibility alias used by existing code.
        self.gestures_sent = self.gesture_timestamps

    def set_embedding(self, embedding: Any) -> None:
        self.embedding = embedding

    def mark_auth_attempt(self, now: float | None = None) -> None:
        self.last_auth_attempt_time = now if now is not None else time.time()

    def can_attempt_auth(self, cooldown_seconds=AUTH_COOLDOWN, now: float | None = None) -> bool:
        now = now if now is not None else time.time()
        if self.last_auth_attempt_time is None:
            return True
        return (now - self.last_auth_attempt_time) > cooldown_seconds

    def authenticate(self, auth_token: str, now: float | None = None) -> None:
        timestamp = now if now is not None else time.time()
        self.last_auth_time = timestamp
        self.last_auth_attempt_time = timestamp
        self.auth_token = auth_token

    def invalidate_auth(self) -> None:
        self.auth_token = None

    def is_authenticated(self, lock_seconds: float | None = None, now: float | None = None) -> bool:
        if self.auth_token is None:
            return False
        if lock_seconds is None or self.last_auth_time is None:
            return True
        now = now if now is not None else time.time()
        return (now - self.last_auth_time) <= lock_seconds

    def can_trigger_gesture(self, gesture_id: int, cooldown_seconds: float = GESTURE_COOLDOWN, now: float | None = None) -> bool:
        now = now if now is not None else time.time()
        last_time = self.gesture_timestamps.get(gesture_id)
        if last_time is None:
            return True
        return (now - last_time) > cooldown_seconds

    def mark_gesture_triggered(self, gesture_id: int, now: float | None = None) -> None:
        self.gesture_timestamps[gesture_id] = now if now is not None else time.time()

    def to_embedding_list(self) -> list[float]:
        if hasattr(self.embedding, "tolist"):
            return self.embedding.tolist()
        return list(self.embedding)

    # Backward-compatible names used by your initial WIP implementation.
    def cooldown_check(self) -> bool:
        return self.can_attempt_auth(AUTH_COOLDOWN)

    def gesture_cooldown_check(self, gesture_id: int) -> bool:
        return self.can_trigger_gesture(gesture_id, GESTURE_COOLDOWN)

    def send_gesture(self, gesture_id: int) -> None:
        self.mark_gesture_triggered(gesture_id)