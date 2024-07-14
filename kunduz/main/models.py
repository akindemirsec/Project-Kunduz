from django.db import models
from django.contrib.auth.models import AbstractUser

# Create your models here.
class User(AbstractUser):
    pass

class Cluster(models.Model):
    name = models.CharField(max_length=255)
    owner = models.ForeignKey(User, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)

class CVE(models.Model):
    cve_id = models.CharField(max_length=50, unique=True)
    description = models.TextField()

class Asset(models.Model):
    name = models.CharField(max_length=255)
    cluster = models.ForeignKey(Cluster, related_name='assets', on_delete=models.CASCADE)

class SBOM(models.Model):
    name = models.CharField(max_length=255)
    cluster = models.ForeignKey(Cluster, related_name='sboms', on_delete=models.CASCADE)

class Alarm(models.Model):
    message = models.TextField()
    cluster = models.ForeignKey(Cluster, related_name='alarms', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)

class CVE(models.Model):
    cve_id = models.CharField(max_length=20, unique=True)
    description = models.TextField()
    published_date = models.DateField()
    last_modified_date = models.DateField()

    def __str__(self):
        return self.cve_id