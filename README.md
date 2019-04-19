# Hodor - Hold the Door
An OAuth2/OpenID Connector External Auth Server for Ambassador

## How does it work

In short, hodor checks every incoming request for a valid token, otherwise it redirects the user to login.

Hodor works in 3 stages: 
- Authentication
- Authorization
- Issue New Credential

## Road map

- [ ] Authentication
  - [ ] simple oauth2 (web flow)
  - [ ] oauth2 with user info
  - [ ] oidc
- [ ] Authorization
- [ ] Issuing
