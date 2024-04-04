import os
from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    DB_URI: str
    ECHO_SQL: bool
    JWT_SECRET_KEY: str
    JWT_EXPIRE_MIN: int

    model_config = SettingsConfigDict(
        env_file=Path(__file__).parent.parent / "settings.env",
        case_sensitive=True,
    )


settings = Settings.model_validate({})