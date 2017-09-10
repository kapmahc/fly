package com.github.kapmahc.fly.nut.models;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@Entity
@Table(name = "policies", indexes = {
        @Index(columnList = "user_id,role_id", unique = true)
})
public class Policy implements Serializable {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @ManyToOne
    @JoinColumn(nullable = false, updatable = false)
    private User user;
    @ManyToOne
    @JoinColumn(nullable = false, updatable = false)
    private Role role;
    @Column(nullable = false, name = "_begin")
    @Temporal(TemporalType.DATE)
    private Date begin;
    @Column(nullable = false, name = "_end")
    @Temporal(TemporalType.DATE)
    private Date end;
    @Column(nullable = false)
    private Date updatedAt;
    @Column(nullable = false, updatable = false)
    private Date createdAt;

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }

    public Date getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(Date updatedAt) {
        this.updatedAt = updatedAt;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public User getUser() {
        return user;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public Role getRole() {
        return role;
    }

    public void setRole(Role role) {
        this.role = role;
    }

    public Date getBegin() {
        return begin;
    }

    public void setBegin(Date begin) {
        this.begin = begin;
    }

    public Date getEnd() {
        return end;
    }

    public void setEnd(Date end) {
        this.end = end;
    }


}
