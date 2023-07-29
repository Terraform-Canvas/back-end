<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
[![Contributors][contributors-shield]][contributors-url] [![Forks][forks-shield]][forks-url] [![Stargazers][stars-shield]][stars-url] [![Issues][issues-shield]][issues-url] [![Pull Request][pr-shield]][pr-url] [![MIT License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/Terraform-Canvas/back-end">
    <img src="./images/canvas-logo.jpg" alt="Logo" width="200" height="200">
  </a>

<h3 align="center">Terraform-Canvas</h3>

  <p align="center">
    Terraform Cloud Infrastructure Provisioning Web Services Project with Visual Programming
    <br />
    <a href="https://facerain.notion.site/e393c21c423e46318f1dd21a3a9ed428?v=cf7ba34920154548a7d0303f27c7710b&pvs=4"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/Terraform-Canvas/back-end/README_kor.md">한국어</a>
    ·
    <a href="https://github.com/Terraform-Canvas/back-end/README.md">English</a>
    <br />
    <br />
    <a href="https://github.com/Terraform-Canvas/back-end/issues">Report Issues</a>
    ·
    <a href="https://github.com/Terraform-Canvas/back-end/pulls">Pull Requests</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#configuration">Configuration</a></li>
      </ul>
    </li>
    <li><a href="#architecture">Architecture</a></li>
    <li><a href="#rest-api">REST API</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## :mag: About The Project
Many companies are choosing Terraform as an IaC tool for transitioning from existing cloud and on-premise environments to cloud native environments. However, they are suffering a lot due to Terraform's high entry barriers. Therefore, we want to lower the barriers to Terraform's entry through "Terraform Cloud Infrastructure Provisioning Web Service with Visual Programming." This enables start-up and SI teams that want to introduce a new cloud-native environment and TF teams that want to test and prototype the IaC environment quickly.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### :card_file_box: Built With
#### :bulb: Language
[![Go][Go]][Go-url]
#### :bulb: Infrastructure
[![Terraform][Terraform]][Terraform-url] [![Kubernetes][Kubernetes]][Kubernetes-url] [![aws][aws]][aws-url] [![OCI][OCI]][OCI-url]
#### :bulb: Environment (CI/CD, Package tools...)
[![Github-actions][Github-actions]][Github-actions-url] [![Helm][Helm]][Helm-url] [![Accordian][Accordian]][Accordian-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## :rocket: Getting Started

### :zap: Prerequisites
Download and install packages and associated dependencies via `go get`
* go
  ```sh
  go get .
  ```

### :pencil2: Configuration
Setting environment variables through `.env`
```env
# .env

# Stage status to start server:
#   - "dev", for start server without graceful shutdown
#   - "prod", for start server with graceful shutdown
STAGE_STATUS="dev"

# Server settings:
SERVER_HOST="0.0.0.0"
SERVER_PORT=8000
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY="secret"
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
JWT_REFRESH_KEY="refresh"
JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT=720

# OCI SDK settings:
tenancyID=tenancy
userID=user
fingerprint=fingerprint
privateKeyFile=filePath
region=us-ashburn-1
compartmentID=compartmentID
privateKeyPass=

# AWS SDK settings:
AWS_ACCESS_KEY=USER_ACCESS_KEY
AWS_SECRET_KEY=USER_SECRET_KEY
AWS_REGION=USER_REGION
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



## :globe_with_meridians: Architecture
### :triangular_flag_on_post: Overall Service Configuration Architecture
![service](./images/service-architecture.png)

### :triangular_flag_on_post: Development Environment Architecture
![env](./images/env-architecture.png)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## :memo: REST API
[Login new](https://www.notion.so/facerain/new-c4285cb8039844eeb4f6ac6fc3db31e0?pvs=4)

[Login refresh](https://www.notion.so/facerain/refresh-5549d45f449c4388b907c4fc03251943?pvs=4)

[Logout](https://www.notion.so/facerain/d72706b87d9f414aa40f57a3bd744bd8?pvs=4)

[Create tf](https://www.notion.so/facerain/tf-60291b66fe524c419f30dc3c13733682?pvs=4)

[Apply tf](https://www.notion.so/facerain/fcff4f41d3ee4b5bb9bcc5fafe180229?pvs=4)

[Save IAM user key](https://www.notion.so/facerain/api-key-e9dc48f44d054aa8929aa976ce7313b8?pvs=4)

[Upload to S3](https://www.notion.so/facerain/S3-27cdcd0c7fdf47a68850e7500db487f6?pvs=4)

[Download to S3](https://www.notion.so/facerain/S3-a45f2ff0d33d465e950cb1b8c159df41?pvs=4)

[Get InstanceTypes](https://www.notion.so/facerain/6f67510b97a34092811c281c737729b1?pvs=4)

[Get AMI](https://www.notion.so/facerain/AMI-9aba2eb13f6842c3b9c91d4240b1f6e2?pvs=4)


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## :fire: Contributing
Please refer to `CONTRIBUTION.txt` for Contribution.

For issues, new functions and requests to modify please follow the following procedure. 🥰

1. Fork the Project
2. Create a Issue when you have new feature or bug, just not Typo fix
3. Create your Feature Branch from dev Branch (`git checkout -b feature/Newfeature`)
4. Commit your Changes (`git commit -m 'feat: add new feature'`)
5. Push to the Branch (`git push origin feature/Newfeature`)
6. Open a Pull Request to dev branch with Issues

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## :closed_lock_with_key: License
Please refer to `LICENSE.txt` for LICENSE.
<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## :speech_balloon: Contact

<table>
  <tbody>
    <tr>
      <td align="center"><a href="https://github.com/Eeap"><img src="https://avatars.githubusercontent.com/u/42088290?v=4" width="100px;" alt=""/><br /><sub><b>Sumin Kim</b></sub></a></td>
      <td align="center"><a href="https://github.com/dusdjhyeon"><img src="https://avatars.githubusercontent.com/u/73868703?v=4" width="100px;" alt=""/><br /><sub><b>Dahyun Kang</b></sub></a></td>
    </tr>
  </tobdy>
</table>

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/Terraform-Canvas/back-end.svg?style=flat
[contributors-url]: https://github.com/Terraform-Canvas/back-end/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Terraform-Canvas/back-end.svg?style=flat
[forks-url]: https://github.com/Terraform-Canvas/back-end/network/members
[stars-shield]: https://img.shields.io/github/stars/Terraform-Canvas/back-end.svg?style=flat
[stars-url]: https://github.com/Terraform-Canvas/back-end/stargazers
[issues-shield]: https://img.shields.io/github/issues/Terraform-Canvas/back-end.svg?style=flat
[issues-url]: https://github.com/Terraform-Canvas/back-end/issues
[pr-url]: https://github.com/Terraform-Canvas/back-end/pulls
[pr-shield]: https://img.shields.io/github/issues-pr/Terraform-Canvas/back-end.svg?style=flat
[license-shield]: https://img.shields.io/github/license/Terraform-Canvas/back-end.svg?style=flat
[license-url]: https://github.com/Terraform-Canvas/back-end/blob/master/LICENSE.txt

[Go]: https://img.shields.io/badge/Go-00ADD8?style=flat&logo=Go&logoColor=white
[Go-url]: https://go.dev/
[Terraform]: https://img.shields.io/badge/Terraform-430098?style=flat&logo=Terraform&logoColor=white
[Terraform-url]: https://www.terraform.io/
[aws]: https://img.shields.io/badge/AmazonAWS-232F3E?style=flat&logo=AmazonAWS&logoColor=white
[aws-url]: https://aws.amazon.com/
[OCI]: https://img.shields.io/badge/Oracle-F80000?style=flat&logo=oracle&logoColor=black
[OCI-url]: https://www.oracle.com/kr/cloud/
[Kubernetes]: https://img.shields.io/badge/Kubernetes-326CE5?style=flat&logo=Kubernetes&logoColor=white
[Kubernetes-url]: https://kubernetes.io/ko/
[Github-actions]: https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat&logo=github-actions&logoColor=white
[Github-actions-url]: https://github.com/features/actions
[Helm]: https://img.shields.io/badge/Helm-326CE5?style=flat&logo=Helm&logoColor=white
[Helm-url]: https://helm.sh/
[Accordian]: https://img.shields.io/badge/Accordian-430098?style=flat&logo=Accordian&logoColor=white
[Accordian-url]: https://accordions.co.kr/
