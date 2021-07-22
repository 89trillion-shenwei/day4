from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)

    @task
    def index1(self):
        self.client.post("/login",{"uid":"1626940055231998000"})

    @task
    def index2(self):
        self.client.post("/register",{})

    @task
    def index3(self):
        self.client.post("/receiveGifts",{"key":"c2906d4b","username":"1626940055231998000"})