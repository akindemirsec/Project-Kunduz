from django.contrib import admin
from .models import User, Cluster, CVE, Asset, SBOM, Alarm

admin.site.register(User)
admin.site.register(Cluster)
admin.site.register(CVE)
admin.site.register(Asset)
admin.site.register(SBOM)
admin.site.register(Alarm)
