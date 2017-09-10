package com.github.kapmahc.fly.nut.models;

import javax.persistence.*;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Entity
@Table(name = "users", indexes = {
        @Index(columnList = "name"),
        @Index(columnList = "email", unique = true),
        @Index(columnList = "uid", unique = true),
        @Index(columnList = "providerType"),
        @Index(columnList = "providerType,providerId", unique = true),
})
public class User implements Serializable {
    public enum Type {
        EMAIL
    }

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @Column(nullable = false)
    private String name;
    @Column(nullable = false, updatable = false)
    private String email;
    @Column(nullable = false, length = 36)
    private String uid;
    private String password;
    @Column(nullable = false)
    private String providerId;
    @Column(nullable = false, length = 16)
    @Enumerated(EnumType.STRING)
    private Type providerType;
    private String logo;
    private long signInCount;
    @Column(length = 45)
    private String lastSignInIp;
    private Date lastSignInAt;
    @Column(length = 45)
    private String currentSignInIp;
    private Date currentSignInAt;
    private Date confirmedAt;
    private Date lockedAt;
    @Column(nullable = false)
    private Date updatedAt;
    @Column(nullable = false, updatable = false)
    private Date createdAt;
    @OneToMany(mappedBy = "user")
    private List<Log> logs;
    @OneToMany(mappedBy = "user")
    private List<Attachment> attachments;
    @OneToMany(mappedBy = "user")
    private List<Policy> policies;


    public User() {
        logs = new ArrayList<>();
        attachments = new ArrayList<>();
        policies = new ArrayList<>();
    }

    public List<Policy> getPolicies() {
        return policies;
    }

    public void setPolicies(List<Policy> policies) {
        this.policies = policies;
    }

    public List<Attachment> getAttachments() {
        return attachments;
    }


    public void setAttachments(List<Attachment> attachments) {
        this.attachments = attachments;
    }

    public List<Log> getLogs() {
        return logs;
    }

    public void setLogs(List<Log> logs) {
        this.logs = logs;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getUid() {
        return uid;
    }

    public void setUid(String uid) {
        this.uid = uid;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getProviderId() {
        return providerId;
    }

    public void setProviderId(String providerId) {
        this.providerId = providerId;
    }

    public Type getProviderType() {
        return providerType;
    }

    public void setProviderType(Type providerType) {
        this.providerType = providerType;
    }

    public String getLogo() {
        return logo;
    }

    public void setLogo(String logo) {
        this.logo = logo;
    }

    public long getSignInCount() {
        return signInCount;
    }

    public void setSignInCount(long signInCount) {
        this.signInCount = signInCount;
    }

    public String getLastSignInIp() {
        return lastSignInIp;
    }

    public void setLastSignInIp(String lastSignInIp) {
        this.lastSignInIp = lastSignInIp;
    }

    public Date getLastSignInAt() {
        return lastSignInAt;
    }

    public void setLastSignInAt(Date lastSignInAt) {
        this.lastSignInAt = lastSignInAt;
    }

    public String getCurrentSignInIp() {
        return currentSignInIp;
    }

    public void setCurrentSignInIp(String currentSignInIp) {
        this.currentSignInIp = currentSignInIp;
    }

    public Date getCurrentSignInAt() {
        return currentSignInAt;
    }

    public void setCurrentSignInAt(Date currentSignInAt) {
        this.currentSignInAt = currentSignInAt;
    }

    public Date getConfirmedAt() {
        return confirmedAt;
    }

    public void setConfirmedAt(Date confirmedAt) {
        this.confirmedAt = confirmedAt;
    }

    public Date getLockedAt() {
        return lockedAt;
    }

    public void setLockedAt(Date lockedAt) {
        this.lockedAt = lockedAt;
    }

    public Date getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(Date updatedAt) {
        this.updatedAt = updatedAt;
    }

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }
}
