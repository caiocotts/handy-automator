from User import User


class UserRegistry:
    def __init__(self):
        self.users: dict[str, User] = {}

    def add_user(self, user_id: str, user: User) -> None:
        self.users[user_id] = user

    def get_user(self, user_id: str) -> User | None:
        return self.users.get(user_id)

    def get_or_create_user(self, user_id: str, embedding=None) -> User:
        user = self.get_user(user_id)
        if user is None:
            user = User(user_id, embedding)
            self.add_user(user_id, user)
        elif embedding is not None:
            user.set_embedding(embedding)
        return user

    def upsert_from_embeddings(self, embeddings: dict[str, object]) -> None:
        for user_id, embedding in embeddings.items():
            self.get_or_create_user(user_id, embedding)

    def has_authenticated_user(self, lock_seconds: float | None = None, now: float | None = None) -> bool:
        for user in self.users.values():
            if user.is_authenticated(lock_seconds=lock_seconds, now=now):
                return True
        return False

    def get_authenticated_user(self, user_id: str, lock_seconds: float | None = None, now: float | None = None) -> User | None:
        user = self.get_user(user_id)
        if user is None:
            return None
        if user.is_authenticated(lock_seconds=lock_seconds, now=now):
            return user
        return None

    def all_users(self) -> list[User]:
        return list(self.users.values())
        