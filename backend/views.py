from datetime import timedelta
from flask import Flask, request, jsonify, send_from_directory, Blueprint, g, redirect, Response
from flask_mail import Message
from flask_cors import CORS
from urllib.parse import unquote
import psycopg2
import os
import re
from flask_jwt_extended import (
    JWTManager,
    create_access_token,
    jwt_required,
    get_jwt_identity,
    decode_token
)
from flask_bcrypt import Bcrypt
import uuid
import stripe
import requests
import json

from database import connect_to_database, get_db

import models.card
import services
import utils

bp = Blueprint("bp", __name__)


def protected(func):
    def mail_wrapper(*args, **kwargs):
        auth_header = request.headers.get('Authorization')
        if not auth_header or auth_header != os.getenv("MAIL_PASSWORD"):
            return jsonify({"error": "Unauthorized"}), 401

        result = func(*args, **kwargs)
        return result
    return mail_wrapper

@bp.route("/api/send", methods=["POST"])
@protected
def send_mail():
    data = request.get_json()
    subject = data.get("subject")
    if subject == "":
        return jsonify({"message": "Email needs a subject"}), 400
    recipient = data.get("recipient")
    if subject == "":
        return jsonify({"message": "Email needs a recipient"}), 400
    body = data.get("body")
    message = Message(subject, recipients=[recipient], body=body)
    g.mail.send(message)
    return jsonify({}), 200


@bp.route("/api/login", methods=["POST"])
def login():
    data = request.get_json()
    params = {
        "email": data.get("email"),
        "password": data.get("password")
    }
    headers = {
        "Authorization": request.headers.get("Authorization"),
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/login/", headers=headers, json=params)
    print(response)
    print(response.text)
    return jsonify(response.json()), response.status_code

def generate_temp_token(user_id):
    expires = timedelta(minutes=5)  # Token expires in 24 hours
    return create_access_token(identity=user_id, expires_delta=expires)

@bp.route("/api/request-reset", methods=["POST"])
def request_password_reset():
    data = request.get_json()
    params = {
        "email": data.get("email")
    }
    email = data.get("email")
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/request-reset/", headers=headers, json=params)
    print(response)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/reset-password", methods=["POST"])
def reset_password():
    data = request.get_json()
    token = data.get("token")
    new_password = data.get("new_password")
    
    params = {
        "token": data.get("token"),
        "new_password": data.get("new_password")
    }
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/reset-password/", headers=headers, json=params)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/email-validate", methods=["GET"])
@jwt_required()
def resend_email_validation():
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/email-validate/", headers=headers)
    print(response)
    print(response.text)
    return jsonify(response.json()), response.status_code
    
    
@bp.route("/api/email-validate", methods=["POST"])
def validate_email():
    data = request.get_json()
    token = data.get("token")
    
    params = {
        "token": data.get("token"),
    }
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/email-validate/", headers=headers, json=params)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/auth", methods=["GET"])
@jwt_required()
def check_token():
    return jsonify({}), 200


@bp.route("/api/cards", methods=["GET"])
@jwt_required()
def get_cards():
    auth_header = request.headers.get("Authorization")
    if not auth_header:
        return jsonify({"error": "Authorization header is missing"}), 401
    
    search_term = request.args.get("search_term", None)
    partial = request.args.get("partial", False)
    sort_method = request.args.get("sort_method", "id")
    inactive = request.args.get("inactive", False)
    
    headers = {
        "Authorization": auth_header
    }
    
    # Forward the request to the Go backend
    params = {
        "search_term": search_term,
        "partial": partial,
        "sort_method": sort_method,
        "inactive": inactive
    }
    print("http://" + os.getenv("FILES_HOST") + "/api/cards/")
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/cards/", params=params, headers=headers)
    print(response)
    
    return jsonify(response.json()), response.status_code

@bp.route("/api/cards/next", methods=["POST"])
@jwt_required()
def generate_next_id():
    card_type = request.json.get("card_type", None)
    params = {
        "card_type": request.json.get("card_type", "")
    }
    
    headers = {
        "Authorization": request.headers.get("Authorization")
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/next/", json=params, headers=headers)
    print(response)
    print(response.text)
    return jsonify(response.json()), response.status_code


@bp.route("/api/cards", methods=["POST"])
@jwt_required()
def create_card():
    card = {
        "title": request.json.get("title"),
        "body": request.json.get("body", ""),
        "card_id": request.json.get("card_id"),
        "link": request.json.get("link", ""),
    }

    headers = {
        "Authorization": request.headers.get("Authorization"),
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/cards/", headers=headers, json=card)
    print(response.text)
    return jsonify(response.json()), response.status_code


@bp.route("/api/cards/<path:id>", methods=["GET"])
@jwt_required()
def get_card(id):

    auth_header = request.headers.get("Authorization")
    if not auth_header:
        return jsonify({"error": "Authorization header is missing"}), 401
    # Forward the request to the Go backend
    headers = {
        "Authorization": auth_header
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/cards/" + str(id) + "/", headers=headers)
    print(response.text)
    return jsonify(response.json()), response.status_code


@bp.route("/api/cards/<path:id>", methods=["PUT"])
@jwt_required()
def update_card(id):
    card = {
        "card_id": request.json.get("card_id"),
        "title": request.json.get("title"),
        "body": request.json.get("body"),
        "link": request.json.get("link"),
    }

    headers = {
        "Authorization": request.headers.get("Authorization"),
        "Content-Type": "application/json"
    }
    response = requests.put("http://" + os.getenv("FILES_HOST") + "/api/cards/" + str(id) + "/", headers=headers, json=card)
    print(response.text)
    return jsonify(response.json()), response.status_code


@bp.route("/api/cards/<path:id>", methods=["DELETE"])
@jwt_required()
def delete_card(id):
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.delete("http://" + os.getenv("FILES_HOST") + "/api/cards/" + str(id) + "/", headers=headers)
    print(response.text)
    if response.status_code != 204: 
        return jsonify({}), response.status_code
    else: 
        return jsonify(response.text), response.status_code
   


@bp.route("/api/users", methods=["POST"])
def create_user():

    headers = {
        "Content-Type": "application/json"
    }
    response = requests.post("http://" + os.getenv("FILES_HOST") + "/api/users/", headers=headers, json=request.json)
    print(response)
    print(response.text)
    return jsonify({}), response.status_code


def admin_only(func):
    def wrapper(*args, **kwargs):

        current_user = get_jwt_identity()  # Extract the user identity from the token
        user = services.query_full_user(current_user)
        if not user["is_admin"]:
            return jsonify({}), 401

        result = func(*args, **kwargs)
        return result
    return wrapper
    
    
@bp.route("/api/users", methods=["GET"])
@jwt_required()
@admin_only
def get_users():
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/users/", headers=headers)
    return jsonify(response.json()), response.status_code

@bp.route("/api/users/<path:id>", methods=["GET"])
@jwt_required()
def get_user(id):
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/users/" + str(id) + "/", headers=headers)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/users/<path:id>/subscription", methods=["GET"])
@jwt_required()
def get_user_subscription(id: int):
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/users/" + str(id) + "/subscription/", headers=headers)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/current", methods=["GET"])
@jwt_required()
def get_current_user():
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/current/", headers=headers)
    print("asdas")
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/users/<path:id>", methods=["PUT"])
@jwt_required()
def update_user(id):
    user = services.query_full_user(id)
    old_email = user["email"]
    
    is_admin = request.json.get("is_admin")
    if "is_admin" not in request.json:
        is_admin = user["is_admin"]
    user = {
        "is_admin": is_admin,
        "username": request.json.get("username"),
        "email": request.json.get("email")
    }
    headers = {
        "Authorization": request.headers.get("Authorization"),
        "Content-Type": "application/json"
    }
    response = requests.put("http://" + os.getenv("FILES_HOST") + "/api/users/" + str(id) + "/", headers=headers, json=user)
    print(response.text)
    return jsonify(response.json()), response.status_code

@bp.route("/api/files/upload", methods=["POST"])
@jwt_required()
def upload_file():
    print("do we reqch here?")
    current_user = get_jwt_identity()  # Extract the user identity from the token

    if "file" not in request.files:
        return jsonify({"error": "No file part"}), 400

    if "card_pk" not in request.form:
        return jsonify({"error": "No PK given"}), 400
        
    file = request.files["file"]
    card_pk = request.form["card_pk"]
    if card_pk == "undefined":
        card_pk = None
    if file.filename == "":
        return jsonify({"error": "No selected file"}), 400

    # Proxy the request to the new backend
    files = {'file': (file.filename, file.stream, file.mimetype)}
    data = {'card_pk': card_pk}
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }

    response = requests.post(f"http://{os.getenv('FILES_HOST')}/api/files/upload/", files=files, data=data, headers=headers)

    print(response)
    print(response.status_code)
    print(response.text)
    if response.status_code != 201:
        return jsonify({"error": response.text}), response.status_code
    print(response.json())

    return jsonify(response.json()), 201

@bp.route("/api/files", methods=["GET"])
@jwt_required()
def get_all_files():
    auth_header = request.headers.get("Authorization")
    if not auth_header:
        return jsonify({"error": "Authorization header is missing"}), 401
    # Forward the request to the Go backend
    headers = {
        "Authorization": auth_header
    }
    response = requests.get("http://" + os.getenv("FILES_HOST") + "/api/files", headers=headers)
    print(response)
    print(response.text)
    #print(response.json())

    # Return the response from the Go backend to the client
    return jsonify(response.json()), response.status_code


@bp.route("/api/files/download/<int:file_id>", methods=["GET"])
@jwt_required()
def download_file(file_id):
    current_user = get_jwt_identity()  # Extract the user identity from the token

    # Check file permissions
    if not services.check_file_permission(file_id, current_user):
        return jsonify({}), 401

    # Proxy the request to the Go backend
    auth_header = request.headers.get("Authorization")
    headers = {
        "Authorization": auth_header
    }

    response = requests.get(f"http://{os.getenv('FILES_HOST')}/api/files/download/{file_id}/", headers=headers, stream=True)
    print(response)
    print(response.status_code)

    if response.status_code != 200:
        return jsonify({"error": response.text}), response.status_code


    # Stream the response content back to the client
    return Response(
        response.iter_content(chunk_size=1024),
        content_type=response.headers['Content-Type'],
      )

@bp.route("/api/files/<int:file_id>", methods=["DELETE"])
@jwt_required()
def delete_file(file_id):

    current_user = get_jwt_identity()  # Extract the user identity from the token
    if not services.check_file_permission(file_id, current_user):
        return jsonify({}), 401

    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.delete(f"http://{os.getenv('FILES_HOST')}/api/files/{file_id}/", headers=headers)
    if response.status_code != 200:
        print(response.text)
        return jsonify({"error": response.text}), response.status_code
    return jsonify({}), 200


@bp.route("/api/files/<int:file_id>", methods=["PATCH"])
@jwt_required()
def edit_file(file_id):
    current_user = get_jwt_identity()  # Extract the user identity from the token

    data = request.get_json()
    print(data)
    if not data:
        return jsonify({"error": "No update data provided"}), 400

    # Forward the request to the Go backend
    headers = {
        "Authorization": request.headers.get("Authorization"),
        "Content-Type": "application/json"
    }
    response = requests.patch(f"http://{os.getenv('FILES_HOST')}/api/files/{file_id}/", headers=headers, json=data)

    if response.status_code != 200:
        print(response.text)
        return jsonify({"error": "Failed to update file metadata"}), response.status_code

    updated_file = response.json()
    
    return jsonify(updated_file), response.status_code

@bp.route("/api/admin", methods=["GET"])
@jwt_required()
def check_admin():
    headers = {
        "Authorization": request.headers.get("Authorization"),
    }
    response = requests.get(f"http://{os.getenv('FILES_HOST')}/api/admin/", headers=headers)
    if response.status_code == 400:
        return jsonify(response.json()), response.status_code
    else:
        return jsonify({}), response.status_code

@bp.route("/api/webhook", methods=["POST"])
def webhook():

    payload = request.data
    sig_header = request.headers["Stripe-Signature"]
    endpoint_secret = g.config["STRIPE_ENDPOINT_SECRET"]

    event = stripe.Webhook.construct_event(
        payload, sig_header, endpoint_secret
    )
    try:
        event = stripe.Webhook.construct_event(
            payload, sig_header, endpoint_secret
        )
    except ValueError as e:
        # Invalid payload
        return jsonify({}), 400
    except stripe.error.SignatureVerificationError as e:
        # Invalid signature
        return jsonify({}), 400
    if event['type'] == 'checkout.session.completed':
        session = stripe.checkout.Session.retrieve(
            event["data"]["object"]["id"],
            expand=["line_items"]
        )
        services.fulfill_subscription(request.json, session)
    elif event['type'] == 'customer.subscription.deleted':
        print(request.json)
    return jsonify({}), 200
    
@bp.route("/api/billing/publishable_key", methods=["GET"])
def get_publishable_key():
    stripe_config = {"publicKey": g.config["STRIPE_PUBLISHABLE_KEY"]}
    return jsonify(stripe_config)

@bp.route("/api/billing/create_checkout_session", methods=['POST'])
@jwt_required()
def create_checkout_session():

    current_user = get_jwt_identity() 
    user = services.query_full_user(current_user)

    
    subscription = services.query_user_subscription(user)
    if subscription["stripe_subscription_status"] == "active":
        return jsonify({"error": "User already has an active subscription"}), 400
    data = request.get_json()
    interval = data.get("interval")
    
    services.sync_stripe_plans()
    plan = services.fetch_plan_information(interval)
    try:
        checkout_session = stripe.checkout.Session.create(
            success_url=g.config['ZETTEL_URL'] + "/app/settings/billing/success?session_id={CHECKOUT_SESSION_ID}",
            cancel_url=g.config['ZETTEL_URL'] + "/app/settings/billing/cancelled",
            payment_method_types=["card"],
            mode="subscription",
            customer_email=user["email"],
            line_items=[
                {
                    'price': plan["stripe_price_id"],
                    'quantity': 1,
                },
            ],
            subscription_data={
                "trial_settings": {"end_behavior": {"missing_payment_method": "cancel"}},
                "trial_period_days": 30,
            }

        )
        g.logger.info("New subscription: %s", user["email"])
        return jsonify({"url": checkout_session.url}), 200
    except Exception as e:
        return jsonify(error=str(e)), 403

@bp.route("/api/billing/success", methods=["GET"])
@jwt_required()
def get_successful_session_data():
    session = stripe.checkout.Session.retrieve(request.args.get('session_id'))
    customer = stripe.Customer.retrieve(session.customer)

    return jsonify(customer), 200
