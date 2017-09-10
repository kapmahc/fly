package com.github.kapmahc.fly.ops.vpn.models;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@Entity
@Table(name = "vpn_logs")
public class Log implements Serializable {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @ManyToOne
    @JoinColumn(nullable = false, updatable = false)
    private User user;
    @Column(length = 45, nullable = false)
    private String trustedIp;
    private int trustedPort;
    @Column(length = 45, nullable = false)
    private String remoteIp;
    private int remotePort;
    private double received;
    private double send;
    @Column(nullable = false, name = "_begin")
    private Date begin;
    @Column(nullable = false, name = "_end")
    private Date end;

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

    public String getTrustedIp() {
        return trustedIp;
    }

    public void setTrustedIp(String trustedIp) {
        this.trustedIp = trustedIp;
    }

    public int getTrustedPort() {
        return trustedPort;
    }

    public void setTrustedPort(int trustedPort) {
        this.trustedPort = trustedPort;
    }

    public String getRemoteIp() {
        return remoteIp;
    }

    public void setRemoteIp(String remoteIp) {
        this.remoteIp = remoteIp;
    }

    public int getRemotePort() {
        return remotePort;
    }

    public void setRemotePort(int remotePort) {
        this.remotePort = remotePort;
    }

    public double getReceived() {
        return received;
    }

    public void setReceived(double received) {
        this.received = received;
    }

    public double getSend() {
        return send;
    }

    public void setSend(double send) {
        this.send = send;
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
