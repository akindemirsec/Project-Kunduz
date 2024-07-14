from django.contrib import admin
from django.urls import path
from main import views

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', views.index, name='index'),
    path('login/', views.login_view, name='login'),
    path('register/', views.register_view, name='register'),
    path('clusters/', views.clusters, name='clusters'),
    path('cves/', views.cves, name='cves'),
    path('alarms/', views.alarms, name='alarms'),
    path('scan/', views.scan, name='scan'),
    path('logout/', views.logout_view, name='logout'),
    path('api/update_progress/', views.update_progress, name='update_progress'),
     path('api/start_cve_update/', views.start_cve_update, name='start_cve_update'),
]
