package com.github.kapmahc.fly.survey.models;

import com.github.kapmahc.fly.nut.models.ContentType;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@Entity
@Table(name = "survey_forms")
public class Form implements Serializable {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @Column(nullable = false)
    private String title;
    @Column(nullable = false)
    @Lob
    private String body;
    @Column(nullable = false, length = 8)
    @Enumerated(EnumType.STRING)
    private ContentType type;
    @Column(nullable = false)
    private Date updatedAt;
    @Column(nullable = false, updatable = false)
    private Date createdAt;
}
