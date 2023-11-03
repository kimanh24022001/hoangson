> auth: no
> play: yes

# Register user
This API allows users to register for an account using either their
email address or phone number. The request body should contain an
email or phone number. The API will validate the value and create a
new account if it is valid.

## Request Body
This is the request body structure.
```JsonReq
UserRegisterReq
```

## Response Body
```JsonRep
UserRegisterRep
```

## Error code
This is the every response if an error is detected.
```Errors
{
    RegisterError_ExistEmail, 
    RegisterError_ExistPhone,
}
```
