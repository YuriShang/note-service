from datetime import datetime
import uuid

from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy import String, ForeignKey, func, event, Boolean
from sqlalchemy.orm import mapped_column, Mapped

from .base import Base


class User(Base):
    __tablename__ = 'users'

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True), nullable=False, 
        unique=True, primary_key=True, default=uuid.uuid4
    )
    username: Mapped[str] = mapped_column(
            String(length=320), unique=True, index=True, nullable=False
        )
    hashed_password: Mapped[str] = mapped_column(String(length=1024), nullable=False)
    register_time: Mapped[datetime] = mapped_column(insert_default=func.now(), nullable=False)
    password_set_time: Mapped[datetime] = mapped_column(insert_default=func.now(), nullable=False)
    is_superuser: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    is_verified: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    

class Note(Base):
    __tablename__ = "notes"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True), nullable=False, 
        unique=True, primary_key=True, default=uuid.uuid4
    )
    user_id: Mapped[uuid.UUID] = mapped_column(ForeignKey("users.id"), nullable=False)
    create_time: Mapped[datetime] = mapped_column(insert_default=func.now())
    text: Mapped[str] = mapped_column(String(128), nullable=False)
    public: Mapped[bool] = mapped_column(nullable=False)

    
@event.listens_for(User.hashed_password, "set", active_history=True)
def receive_modified(target, value, old, initiator):
    """Триггер для обновления времени установки пароля при его изменении"""
    target.password_set_time = datetime.now()
