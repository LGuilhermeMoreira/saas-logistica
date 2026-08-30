package authz

# Por padrão, bloqueia tudo
default allow = false

# admin
allow if {
    role := data.roles[_]
    role.name == input.user.role
    
    perm := role.permissions[_]
    perm.action == "*"
    perm.path == "*"
}

# resto dos cargos
allow if {
    role := data.roles[_]
    role.name == input.user.role
    
    perm := role.permissions[_]
    perm.action == input.method
    perm.path == input.path
}