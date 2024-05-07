import logging
from typing import AsyncGenerator, Optional

from fastapi import Depends
from sqlalchemy import func, select
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from fastapi_users.db import SQLAlchemyUserDatabase

from app.models.models import User
from app.config import settings

logger = logging.getLogger(__name__)


class UserDatabase(SQLAlchemyUserDatabase):
    async def get_by_username(self, username: str) -> Optional[User]:
        statement = select(self.user_table).where(
            func.lower(self.user_table.username) == func.lower(username)
        )
        return await self._get_user(statement)


async_engine = create_async_engine(
    settings.DB_URI,
    pool_pre_ping=True,
    echo=settings.ECHO_SQL,
)

async_session_maker = async_sessionmaker(async_engine, expire_on_commit=False)


async def get_async_session() -> AsyncGenerator[AsyncSession, None]:
    async with async_session_maker() as session:
        try:
            yield session
        except SQLAlchemyError as e:
            logger.exception(e)


async def get_user_db(session: AsyncSession = Depends(get_async_session)):
    yield UserDatabase(session, User)