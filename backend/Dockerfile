FROM python:3
ENV PYTHONUNBUFFERED 1
RUN mkdir /app/
RUN mkdir /app/backend
WORKDIR /app/backend
COPY requirements.txt /app/backend
RUN pip install -r requirements.txt
COPY . /app/backend

EXPOSE 8000
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
