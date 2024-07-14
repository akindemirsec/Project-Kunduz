import os
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.decorators import login_required
from django.shortcuts import render, redirect
from django.utils.dateparse import parse_date
from django.http import JsonResponse
import requests
import csv
import time
from .models import *
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django.utils.dateparse import parse_date
from .models import CVE

from .models import User
from .forms import RegisterForm
from django.contrib import messages

downloaded_size = 0

def index(request):
    return render(request, 'index.html')

def login_view(request):
    if request.method == 'POST':
        username = request.POST['username']
        password = request.POST['password']
        user = authenticate(request, username=username, password=password)
        if user is not None:
            login(request, user)
            return redirect('index')
        else:
            return render(request, 'login.html', {'error': 'Invalid username or password'})
    else:
        return render(request, 'login.html')

def register_view(request):
    if request.method == 'POST':
        username = request.POST['username']
        password = request.POST['password']
        password2 = request.POST['password2']
        
        if password == password2:
            if User.objects.filter(username=username).exists():
                messages.error(request, 'Username already exists')
            else:
                user = User.objects.create_user(username=username, password=password)
                user.save()
                messages.success(request, 'Account created successfully')
                return redirect('login')
        else:
            messages.error(request, 'Passwords do not match')
    
    return render(request, 'register.html')

@login_required
def logout_view(request):
    logout(request)
    return redirect('login')

@login_required
def clusters(request):
    return render(request, 'clusters.html')

@login_required
def cves(request):
    return render(request, 'cves.html')

@login_required
def alarms(request):
    return render(request, 'alarms.html')

@login_required
def scan(request):
    return render(request, 'scan.html')

def start_cve_update(request):
    try:
        url = "https://cve.mitre.org/data/downloads/allitems.csv"
        local_filename = "cve_data.csv"
        response = requests.get(url, stream=True)
        total_length = response.headers.get('content-length')

        if total_length is None:
            return JsonResponse({'error': 'Failed to download CVE data'}, status=500)

        total_length = int(total_length)
        downloaded = 0

        with open(local_filename, 'wb') as f:
            for data in response.iter_content(chunk_size=4096):
                f.write(data)
                downloaded += len(data)
                # Send progress to WebSocket
                progress = int(100 * downloaded / total_length)
                channel_layer = get_channel_layer()
                async_to_sync(channel_layer.group_send)(
                    'progress_group', {
                        'type': 'send_progress',
                        'progress': progress,
                    }
                )

        # Process CSV file
        with open(local_filename, newline='') as csvfile:
            reader = csv.reader(csvfile)
            for row in reader:
                # Insert or update the CVE record in your database
                pass

        os.remove(local_filename)
        return JsonResponse({'status': 'CVE update started'})
    except Exception as e:
        return JsonResponse({'error': f'Failed to start CVE update: {str(e)}'}, status=500)     

@login_required
def start_cve_update(request):
    try:
        update_cve_database()
        return JsonResponse({'status': 'success'})
    except Exception as e:
        return JsonResponse({'status': 'error', 'message': str(e)})

def update_progress(request):
    global downloaded_size
    return JsonResponse({'progress': downloaded_size})

def process_sbom(file_path):
    # Dummy function to process SBOM and generate alarms
    # Implement your logic to process the SBOM file and generate alarms
    pass
