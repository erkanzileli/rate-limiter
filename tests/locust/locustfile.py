import time
from locust import HttpUser, task, between

class QuickstartUser(HttpUser):
    wait_time = between(1, 2.5)

    @task
    def unlimited_endpoint(self):
        self.client.get("/packages")
#
#     @task
#     def limited_endpoint(self):
#         self.client.get("/xyz")
