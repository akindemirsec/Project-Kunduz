from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('login/', views.login_view, name='login'),
    path('register/', views.register_view, name='register'),
    path('logout/', views.logout_view, name='logout'),
    path('clusters/', views.clusters, name='clusters'),
    path('cves/', views.cves, name='cves'),
    path('alarms/', views.alarms, name='alarms'),
    path('scan/', views.scan, name='scan'),
    path('api/update_progress/', views.update_progress, name='update_progress'),
    path('api/update_cve_database/', views.update_cve_database, name='update_cve_database'),
]
