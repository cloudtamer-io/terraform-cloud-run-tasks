# Terraform Cloud Integration with Run Tasks

Terraform Cloud is a managed service by HashiCorp that makes it easy to run Terraform in production. With the release of run tasks for Terraform Cloud, we are able to easily integrate into the workflow of a Terraform deployment by executing webhooks after the plan stage, but before the apply stage. The tasks can be set to either advisory or mandatory to show errors as informational or fail the pipeline, respectively.

We have two event hooks available:

- Savings Opportunities - shows the monthly forecast on a project as well as the cost savings available.
- Compliance - shows the number of compliance findings on a project and will throw an error if there are any critical findings.

This repository contains the code to launch an EC2 instance to run the webhook.

You can access the deployment guide [here](https://cloudtamer.zendesk.com/hc/en-us/articles/4408728893325).
